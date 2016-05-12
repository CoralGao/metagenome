import csv 
import sys

alt_map = {'ins':'0'}
complement = {'A': 'T', 'C': 'G', 'G': 'C', 'T': 'A'} 

def reverse_complement(seq):    
    for k,v in alt_map.iteritems():
        seq = seq.replace(k,v)
    bases = list(seq) 
    bases = reversed([complement.get(base,base) for base in bases])
    bases = ''.join(bases)
    for k,v in alt_map.iteritems():
        bases = bases.replace(v,k)
    return bases

infile = sys.argv[1]
outfile = sys.argv[2]

f = open(infile,'r')
out = open(outfile, 'w')

out.write(f.readline())

l = 60
seq = ""
for line in f.readlines():
	out.write(line)
	line = line.strip()
	seq += line

rev = reverse_complement(seq)

i = 0
print len(rev)
while i < len(rev):
    e = min(len(rev), i+l)
    out.write(rev[i:i+l])
    out.write("\n")
    i = e
out.close()
