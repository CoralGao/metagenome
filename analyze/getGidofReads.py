import sys

titlefile = sys.argv[1]
readsfile = sys.argv[2]

title = open(titlefile,"r")
reads = open(readsfile,"r")

titles = []
for line in title.readlines():
	line = line.strip()
	titles.append(line)
title.close()

i = 0
print("realGid")
for line in reads.readlines():
	if i%2 == 0:
		line= line.strip()
		for s in range(len(titles)):
			if titles[s] in line:
				print(s)
	i += 1
reads.close()