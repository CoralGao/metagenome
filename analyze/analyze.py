import sys

file1 = sys.argv[1]
file2 = sys.argv[2]

realgid = open(file1,"r")
raw = open(file2, "r")

i = -1
for l1 in realgid.readlines():
	l2 = raw.readline().strip()
	l1 = l1.strip()
	# true negative
	# if l1 != l2 and l2!= '-1':
	# 	print i
	# false negative
	if l2 == '-1':
		print i
	i += 1

realgid.close()
raw.close()

# file1 = sys.argv[1]
# analysis = open(file1, "r")

# analysis.readline()

# i = 0
# for line in analysis.readlines():
# 	elements = line.strip().split()
# 	if elements[0] == '-1':
# 		print i
# 	i += 1
# analysis.close()

# for line in analysis.readlines():
# 	elements = line.strip().split()
# 	if elements[0] == elements[1]:
# 		pp.write(line)
# 	elif elements[0] != elements[1] and elements[0] != '-1':
# 		pn.write(line)
# 	elif elements[0] != elements[1] and elements[0] == '-1':
# 		np.write(line)
# analysis.close()

# file1 = sys.argv[1]
# analysis = open(file1, "r")
# res = 0

# for line in analysis.readlines():
# 	elements = line.strip().split()
# 	# if elements[1] == "-1":
# 	if elements[0] == elements[2] and elements[1] != '-1':

# 	# if elements[1] != "-1" and elements[0] != '-1':
# 		res += 1
# print res
# analysis.close()
