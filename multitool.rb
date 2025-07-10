class Multitool < Formula
  desc "Multicloud CLI for object storage, K8s, and more (by tronicum)"
  homepage "https://github.com/tronicum/punchbag-cube-testsuite"
  url "https://github.com/tronicum/punchbag-cube-testsuite/releases/download/v0.1.0/mt-macos-amd64.tar.gz"
  sha256 "<FILL_ME_IN_AFTER_UPLOAD>"
  license "AGPL-3.0-only"

  def install
    bin.install "mt"
  end

  test do
    system "#{bin}/mt", "--help"
  end
end
