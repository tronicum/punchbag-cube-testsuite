#!/usr/bin/env perl
use strict;
use warnings;
use File::Basename;
use Cwd 'abs_path';

my $module_prefix = 'github.com/tronicum/punchbag-cube-testsuite';
my $root = abs_path('.');

sub fix_module_line {
    my ($lines_ref, $file) = @_;
    my @lines = @$lines_ref;
    my $fixed = 0;
    # If first line is not a valid module line, reconstruct it from path
    if (!@lines || $lines[0] !~ /^module\s+\S+\s*$/) {
        my $mod_path = $file;
        $mod_path =~ s/^\Q$root\E\/?//;
        $mod_path =~ s/\/go\.mod$//;
        $mod_path =~ s{^/?}{};
        my $module = $module_prefix;
        if ($mod_path ne '' && $mod_path ne 'go.mod') {
            $module .= "/$mod_path";
        }
        my $rest = '';
        if (@lines && $lines[0] =~ /^(?:module\s+\S+)?(.*)$/) {
            $rest = $1;
        }
        $lines[0] = "module $module\n";
        if ($rest =~ /\S/) {
            my @parts = split(/(?=(?:go\s+\d+(?:\.\d+)*|require\s*\(|replace\s+))/, $rest);
            @parts = map { s/^\s+//; $_ } @parts;
            splice(@lines, 1, 0, map { "$_\n" } grep { /\S/ } @parts);
        }
        $fixed = 1;
    }
    return ($fixed, \@lines);
}

sub process_file {
    my ($file) = @_;
    open my $in, '<', $file or die "Can't read $file: $!";
    my @lines = <$in>;
    close $in;

    my $changed = 0;
    my @newlines;
    my $in_require = 0;

    my ($fixed, $fixed_lines) = fix_module_line(\@lines, $file);
    @lines = @$fixed_lines;
    $changed ||= $fixed;

    foreach my $line (@lines) {
        # Track require block
        if ($line =~ /^\s*require\s*\(\s*$/) {
            $in_require = 1;
            push @newlines, $line;
            next;
        }
        if ($in_require && $line =~ /^\s*\)\s*$/) {
            $in_require = 0;
            push @newlines, $line;
            next;
        }
        # Remove any require/replace for any github.com/tronicum/punchbag-cube-testsuite/ path
        if ($line =~ m{^\s*(require|replace)?\s*github\.com/tronicum/punchbag-cube-testsuite/.*}) {
            $changed = 1;
            next;
        }
        # Remove any line in require block for any github.com/tronicum/punchbag-cube-testsuite/ path
        if ($in_require && $line =~ m{^\s*github\.com/tronicum/punchbag-cube-testsuite/.*}) {
            $changed = 1;
            next;
        }
        push @newlines, $line;
    }

    # Remove empty require blocks
    if ($changed) {
        my $content = join('', @newlines);
        $content =~ s/require\s*\(\s*\)//gs;
        @newlines = split(/\n/, $content, -1);

        open my $out, '>', $file or die "Can't write $file: $!";
        print $out @newlines;
        close $out;
        my $dir = dirname($file);
        print "Fixed: $file\n";
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
        if ($line =~ m{^\s*replace\s+github\.com/tronicum/punchbag-cube-testsuite/.*$}) {
            $changed = 1;
            next;
        }
        push @newlines, $line;
    }
    if ($changed) {
        open my $out, '>', $file or die "Can't write $file: $!";
        print $out @newlines;
        close $out;
        print "Fixed: $file\n";
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