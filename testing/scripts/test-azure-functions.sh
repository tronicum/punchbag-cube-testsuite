#!/bin/bash
# filepath: /Users/ssels/workspace:vs-code:base/punchbag-cube-testsuite/punchbag-cube-testsuite/scripts/test-azure-functions.sh

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
WERFTY_DIR="$PROJECT_ROOT/werfty"

# Default values
API_URL="http://localhost:8080"
SIMULATION_MODE=true
RESOURCE_GROUP="test-rg-$(date +%s)"
LOCATION="eastus"

print_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  --direct          Use direct mode instead of simulation"
    echo "  --api-url URL     API base URL (default: $API_URL)"
    echo "  --resource-group  Azure resource group name"
    echo "  --location        Azure location (default: $LOCATION)"
    echo "  --help           Show this help message"
}

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --direct)
            SIMULATION_MODE=false
            shift
            ;;
        --api-url)
            API_URL="$2"
            shift 2
            ;;
        --resource-group)
            RESOURCE_GROUP="$2"
            shift 2
            ;;
        --location)
            LOCATION="$2"
            shift 2
            ;;
        --help)
            print_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Build werfty if needed
build_werfty() {
    log_info "Building werfty..."
    cd "$WERFTY_DIR"
    if ! go build -o werfty ./cmd; then
        log_error "Failed to build werfty"
        exit 1
    fi
    log_success "werfty built successfully"
}

# Test cube-server connectivity
test_connectivity() {
    log_info "Testing connectivity to $API_URL..."
    if curl -s --max-time 5 "$API_URL/healthz" >/dev/null 2>&1; then
        log_success "Cube-server is reachable"
    else
        log_warning "Cube-server at $API_URL is not reachable"
        if $SIMULATION_MODE; then
            log_error "Simulation mode requires cube-server to be running"
            exit 1
        fi
    fi
}

# Test Azure Monitor creation
test_azure_monitor() {
    log_info "Testing Azure Monitor creation..."
    
    local cmd_args=(
        "provider" "azure" "monitor" "create"
        "--api-url" "$API_URL"
        "--resource-group" "$RESOURCE_GROUP"
        "--workspace-name" "test-workspace-$(date +%s)"
        "--location" "$LOCATION"
        "--output" "json"
    )
    
    if $SIMULATION_MODE; then
        cmd_args+=("--simulation")
    fi
    
    if ./werfty "${cmd_args[@]}"; then
        log_success "Azure Monitor test passed"
    else
        log_error "Azure Monitor test failed"
        return 1
    fi
}

# Test Azure Budget creation
test_azure_budget() {
    log_info "Testing Azure Budget creation..."
    
    local cmd_args=(
        "provider" "azure" "budget" "create"
        "--api-url" "$API_URL"
        "--name" "test-budget-$(date +%s)"
        "--amount" "1000"
        "--resource-group" "$RESOURCE_GROUP"
        "--time-grain" "Monthly"
        "--output" "json"
    )
    
    if $SIMULATION_MODE; then
        cmd_args+=("--simulation")
    fi
    
    if ./werfty "${cmd_args[@]}"; then
        log_success "Azure Budget test passed"
    else
        log_error "Azure Budget test failed"
        return 1
    fi
}

# Test Azure AKS creation
test_azure_aks() {
    log_info "Testing Azure AKS creation..."
    
    local cmd_args=(
        "provider" "azure" "aks" "create"
        "--api-url" "$API_URL"
        "--name" "test-aks-$(date +%s)"
        "--resource-group" "$RESOURCE_GROUP"
        "--location" "$LOCATION"
        "--node-count" "2"
        "--output" "json"
    )
    
    if $SIMULATION_MODE; then
        cmd_args+=("--simulation")
    fi
    
    if ./werfty "${cmd_args[@]}"; then
        log_success "Azure AKS test passed"
    else
        log_error "Azure AKS test failed"
        return 1
    fi
}

# Test provider info
test_provider_info() {
    log_info "Testing provider info..."
    
    if ./werfty provider info azure --api-url "$API_URL" --output table; then
        log_success "Provider info test passed"
    else
        log_error "Provider info test failed"
        return 1
    fi
}

# Main execution
main() {
    log_info "Starting Azure functions test"
    log_info "Mode: $(if $SIMULATION_MODE; then echo 'Simulation'; else echo 'Direct'; fi)"
    log_info "API URL: $API_URL"
    log_info "Resource Group: $RESOURCE_GROUP"
    log_info "Location: $LOCATION"
    
    # Build and test
    build_werfty
    test_connectivity
    
    # Run tests
    local failed_tests=0
    
    test_provider_info || ((failed_tests++))
    test_azure_monitor || ((failed_tests++))
    test_azure_budget || ((failed_tests++))
    test_azure_aks || ((failed_tests++))
    
    # Summary
    echo
    if [[ $failed_tests -eq 0 ]]; then
        log_success "All Azure function tests passed! ✅"
    else
        log_error "$failed_tests test(s) failed ❌"
        exit 1
    fi
}

# Execute main function
main "$@"
