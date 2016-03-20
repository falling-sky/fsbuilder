#! /usr/bin/perl

foreach my $file (@ARGV) {
    fix($file);
}

sub fix {
    my $file = shift @_;
    my $block;

    return unless ( $file =~ m#\.(html|inc|js)$# );

    open( FILE, "<$file" ) or die "$file: $!";
    read FILE, $block, 1000000;
    close FILE;

    @groups = $block =~ m#\{\{(.*?)\}\}#msg;

    @groups = grep( /\\['"]/, @groups );
    return unless ( scalar @groups );

    printf "file $file groups %i\n", scalar @groups;
    $newblock = $block;
    foreach $before (@groups) {
    $DB::single=1;
        $after = $before;
        $after =~ s#\\'#'#g;
        $after =~ s#\\"#"#g;
        die if ($after eq $before);
        
        my $index  = index( $newblock, $before, 0);
        if ($index >= 0) {
          print "Replacing with $after\n";
          substr($newblock,$index,length($before),$after);
        } else {
           die;
        }
    }
    die "wtf" if $newblock eq $block;

    open( FILE, ">$file" ) or die "$file: $!";
    print FILE $newblock;
    close FILE;

}
