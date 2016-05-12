# metagenome
Metagenome Project
0. preprocess reads into one line each:
		python process_reads.py taxon-Empirical.8eebc78e.fna 2x-Empirical.8eebc78e.fna

3. reverse complement all genomes:
	python reverse_com.py 18171.AAYI02000001-AAYI02000004.nuc.fsa     both.18171.AAYI02000001-AAYI02000004.nuc.fsa

1. concatenate all sequences togather, at the same time record the lengths to count the coverage

	 cat both* > all.fasta
	 go run recordPos.go all.fasta 

2. format of paired-end reads:
	grep -A 1 "\.1 " test.fna > reads_1.fna
	grep -A 1 "\.2 " test.fna > reads_2.fna

analysis:
1. get seq & count
	go run guess_sequence.go allboth.fasta single-end/2x-454.8eb42430.fna > raw.txt

2. get real position of reads:
	grep ">" all.fasta > title.txt
	python getGidofReads.py title.txt ../single-end/2x-454.8eb42430.fna > realGid.txt

3. concatenate raw.txt & realGid 
	python analyze.py

generate wrong reads:
1.	get wrong reads id:
	python analyze.py analysis.txt > wrongID.txt

2.	
	python generateReadsWrong.py wrongID.txt 1-9-random-Empirical.b34be078.fna > out.txt

metaphlan.py
1.	build index use bowtie2
	bowtie2-2.0.2/./bowtie2-build all.fasta index/all.fasta
1.0	preprocess reads to 4 lines
	python make_real_read.py 1-9-random-Empirical.b34be078.fna 1-9-random-Empirical.b34be078.fastq
2.	run metaphlan
	python2.7 metaphlan.py 1-10x-454.ad3795d.fastq --nproc 5 --bowtie2db bowtie2db/all.fasta --bt2_ps sensitive-local --bowtie2out metagenome.bt2out.txt
3.	analyze.py
	python2.7 analyze.py metagenome.bt2out.txt 1-10x-454.ad3795d.csv

Focus
1.	build db
		python2.7 dataConfig.py 50genomes/ myDB	
2.	export LD_LIBRARY_PATH="/usr/local/lib/"

3.	In order to insert data into the FOCUS database, you have to run focus only with -d parameter	
		python2.7 focus.py -d myDB 
4.	

ART
1. generate paired-reads
	./art_illumina -sam -i allsingle.fasta -p -l 100 -ss HS25 -f 5 -m 200 -s 10 -o paired_dat
2. concatenate
	python2.7 preprocess.py paired_dat1.fq paired_dat2.fq > paired_reads_5x.fq
