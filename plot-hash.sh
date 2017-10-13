gnuplot -persist <<EOF
set terminal postscript enhanced color
set output "hash-db.eps"
set size 1.5, 1

set multiplot
set origin 0, 0
set title "Hashing Performance" offset char -1, char -20.5 font "Sans, 14.5"
set size 0.55,0.6
set ylabel "hash computation time (us)" offset 3.2,-.2 font "Bold,14"
set grid
set xlabel "key size (bytes) \n\n" offset 0,1 font "Bold,14"
set xrange [1:2000]
set yrange [0:50]
set key at graph .5, 1 vertical samplen 2 font "Sans-Bold,11" spacing 1
set ytic font "Bold, 10.7" offset .6,0 nomirror
set xtic font "Bold, 10.7" offset 0,.5 nomirror
plot "<awk '{print \$3, \$5}' log-hash-test.txt" using 1:(\$2/1000.0) with linespoint axes x1y1 lt 1 lw 3 lc rgb "red" t "hash"

set origin 0.6, 0
set title "golevelDB Performance" offset char -1, char -20.5 font "Sans, 14.5"
set size 0.55,0.6
set ylabel "read time (us)" offset 3.2,-.2 font "Bold,14"
set grid
set xlabel "key size (bytes) \n\n" offset 0,1 font "Bold,14"
set xrange [1:2000]
set yrange [0:2500]
set key at graph .5, 1 vertical samplen 2 font "Sans-Bold,11" spacing 1
set ytic font "Bold, 10.7" offset .6,0 nomirror
set xtic font "Bold, 10.7" offset 0,.5 nomirror
plot "<awk 'BEGIN {nr=1; incr=3} {if (nr == NR) {print \$3, \$5; nr += incr}}' log-db-test.txt" using 1:(\$2/1000.0) with linespoint axes x1y1 lt 1 lw 3 lc rgb "red" t "dbRead"

EOF

epstopdf hash-db.eps
pdf270 hash-db.pdf
rm hash-db.eps
rm hash-db.pdf
mv hash-db-rotated270.pdf hash-db.pdf
pdfcrop hash-db.pdf
mv hash-db-crop.pdf hash-db.pdf

<<NET
set origin 0.37, 0
set title "(b) Image Replication" offset char -1, char -17 font "Sans, 14.5"
set size 0.4,0.5
set ylabel "% Images Placed" offset 3.2,-.2 font "Bold,14"
set xlabel "#Replica\n\n" offset 0,.25 font "Bold,15"
set xrange [1:7]
set yrange [0:100]
set key at graph 1, .4 horizontal samplen 2 font "Sans-Bold,11.8" spacing 1
set ytic font "Bold, 10.7" offset .6,0 nomirror
plot "images/DSC-DSC.txt" using 1:(100-\$2) with linespoint lt 1 lw 2 lc rgb "dark-green" t "LL",\
    "images/ASC-DSC.txt" using 1:(100-\$2) with linespoint lt 2 lw 2 lc rgb "dark-blue" t "SL",\
    "images/DSC-ASC.txt" using 1:(100-\$2) with linespoint lt 3 lw 2 lc rgb "violet" t "LS",\
    "images/ASC-ASC.txt" using 1:(100-\$2) with linespoint lt 4 lw 2 lc rgb "black" t "SS"




unset grid
set origin 0.49, 0
set title "(c) Network Traffic" offset char -1, char -12.1 font "Sans, 13.5"
set size 0.3,0.38
set ylabel "CDF" offset 3.7,0 font "Bold,11"
set xlabel "Reduction (MB)\n\n" offset 0,1.3 font "Bold,11"
set xrange [0:2000]
set yrange [0:1]
set key at graph 1,0.5 horizontal samplen 1.5 font "Sans-Bold,10.8" spacing 1
set xtic 0, 500 font "Bold, 10.7" offset 0,.5 nomirror
set ytic font "Bold, 10.7" offset .6,0
plot "cloudnet/cdf-absolute-net.txt" using 1:2 with lines lt 1 lw 2 lc rgb "black" title "cloudnet"  axis x1y1,\
     "progress/cdf-absolute-net.txt" using 1:2 with lines lt 2 lw 2 lc rgb "dark-red" title "progress-10\%"  axis x1y1,\
     "progress-30/cdf-absolute-net.txt" using 1:2 with lines lt 3 lw 2 lc rgb "dark-blue" title "progress-30\%"  axis x1y1,\
     "progress-50/cdf-absolute-net.txt" using 1:2 with lines lt 4 lw 2 lc rgb "dark-violet" title "progress-50\%"  axis x1y1
NET
