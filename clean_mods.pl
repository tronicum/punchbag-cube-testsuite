#!/usr/bin/env perl
use strict;
use warnings;
use File::Basename;
use Cwd 'abs_path';

my $root = abs_path('/Users/ssels/workspace:vs-code:base/punchbag-cube-testsuite');

# List all your local modules here (relative to your repo root)
my @local_modules = qw(
    shared
    store
    multitool
    generator
    punchbag
    server
    werfty
    werfty-transformator
    terraform-multicloud-provider
);

my $local_regex = join('|', map { quotemeta($_) } @local_modules);

sub process_file {
    my ($file) = @_;
    my $abs_file = abs_path($file);
    my $mod_path = $abs_file;
    $mod_path =~ s/\Q$root\E\/*//;
    $mod_path =~ s{/go\.mod$}{};
    $mod_path =~ s{/$}{};

    my $module_line = "module github.com/tronicum/punchbag-cube-testsuite";
    if ($mod_path ne '' && $mod_path ne 'go.mod') {
        $module_line .= "/$mod_path";
    }
    $module_line .= "\n";

    open my $in,  '<', $file or die "Can't read $file: $!";
    my @lines = <$in>;
    close $in;

    my $changed = 0;
    my @newlines;

    # Always set the correct module line
    if (!defined $lines[0] || $lines[0] !~ /^module\s+github\.com\/tronicum\/punchbag-cube-testsuite/) {
        $lines[0] = $module_line;
        $changed = 1;
    } elsif ($lines[0] ne $module_line) {
        $lines[0] = $module_line;
        $changed = 1;
    }

    foreach my $i (1..$#lines) {
        my $line = $lines[$i];
        # Remove any require line for any local module (block or single-line)
        if ($line =~ m{^\s*(?:require\s+)?github\.com/tronicum/punchbag-cube-testsuite/($local_regex)\b.*$}) {
            $changed = 1;
            next;
        }
        # Remove any replace line for any local module
        if ($line =~ m{^\s*replace\s+github\.com/tronicum/punchbag-cube-testsuite/($local_regex)\b.*$}) {
            $changed = 1;
            next;
        }
        push @newlines, $line;
    }

    # Prepend the module line
    unshift @newlines, $lines[0];

    # Remove empty require blocks
    if ($changed) {
        my $content = join('', @newlines);
        $content =~ s/require\s*\(\s*\)//gs;
        @newlines = split(/\n/, $content, -1);

        open my $out, '>', $file or die "Can't write $file: $!";
        print $out @newlines;
        close $out;
        my $dir = dirname($file);
        system("cd '$dir' && go mod tidy");
    }
}

sub process_work_file {
    my ($file) = @_;
    open my $in, '<', $file or die "Can't read $file: $!";
    my @lines = <$in>;
    close $in;
    my $changed = 0;
    my @newlines;
    foreach my $line (@lines) {
        # Remove any replace line for any local module
        if ($line =~ m{^\s*replace\s+github\.com/tronicum/punchbag-cube-testsuite/($local_regex)\b.*$}) {
            $changed = 1;
            next;
        }
        push @newlines, $line;
    }
    if ($changed) {
        open my $out, '>', $file or die "Can't write $file: $!";
        print $out @newlines;
        close $out;
    }
}

sub find_and_fix {
    my ($dir) = @_;
    opendir(my $dh, $dir) or die "Can't open $dir: $!";
    while (my $entry = readdir($dh)) {
        next if $entry eq '.' or $entry eq '..';
        my $path = "$dir/$entry";
        if (-d $path) {
            find_and_fix($path);
        } elsif ($entry eq 'go.mod') {
            process_file($path);
        } elsif ($entry eq 'go.work') {
            process_work_file($path);
        }
    }
    closedir($dh);
}

find_and_fix($root);