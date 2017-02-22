import os
import shutil


if os.path.exists("release") == False:
	os.mkdir("release")

# code = os.system("go build -o ./release/chatserver_dawin main.go")
code = os.system("GOOS=linux GOARCH=amd64 go build -o ./release/epubOnline_linux64 main.go")
if code > 0:
	raise ValueError("build error")

dirSrc = "./"
dirDest = "./release"

def getDestPath(path, src, dest):
	if src == "./":
		if path.find("./") == 0:
			path = path[2:]
		return os.path.join(dest, path)

	return path.replace(src, dest)

def hasHidePath(path):
	paths = path.split(os.sep)
	if paths[0] == ".":
		paths = paths[1:]

	if len(paths[0]) <= 0:
		return False

	for p in paths:
		if p[0] == ".":
			return True

	return False

# set files and paths that should be exclued from release files
ignoredDirs = ["controller", "tests", "release", "test", "epub"]
ignoredFiles = ["build_release.py", "main.go", "conf/config.json"]


ignoredPaths = []

for d in ignoredDirs:
	ignoredPaths.append(os.path.join(dirSrc, d))

for f in ignoredFiles:
	ignoredPaths.append(os.path.join(dirSrc, f))



def sha1OfFile(filepath):
    import hashlib
    sha = hashlib.sha1()
    with open(filepath, 'rb') as f:
        while True:
            block = f.read(2**10) # Magic number: one-megabyte blocks.
            if not block: break
            sha.update(block)
        return sha.hexdigest()

def isPathIgnored(filepath):
	for p in ignoredPaths:
		index = filepath.find(p)
		if index == 0:
			return True

	return False

def buildRelease():
	print("bulid release tool start...")

	for root, dirs, files in os.walk(dirSrc):
		if hasHidePath(root):
			continue


		for x in dirs:
			srcPath = os.path.join(root,x)
			destRelativePath = getDestPath(srcPath, dirSrc, dirDest)

			if hasHidePath(x):
				continue

			if isPathIgnored(srcPath) == True:
				continue

			if os.path.exists(destRelativePath) == False:
				try:
					os.mkdir(destRelativePath)
					print("duplicate dir ", srcPath," OK")
				except OSError as e:
					print("mkdir error: ", e)
					raise e


		for x in files:
			filePath = os.path.join(root,x)
			destRelativePath = getDestPath(filePath, dirSrc, dirDest)

			if hasHidePath(x):
				continue

			if isPathIgnored(filePath) == True:
				continue

			if os.path.exists(destRelativePath) == False:
				try:
					shutil.copy(filePath, destRelativePath)
					print("copy file ", filePath, " OK")
				except Exception as e:
					print("copy file error: ", filePath, " ", e)
					raise e
			else:
				destSha1 = sha1OfFile(destRelativePath)
				srcSha1 = sha1OfFile(filePath)
				if destSha1 != srcSha1:
					try:
						os.remove(destRelativePath)
						shutil.copy(filePath, destRelativePath)
						print("copy file ", filePath, " OK")
					except Exception as e:
						raise e

	print("build OK")	
	print("----------------------------------------------------------------")

buildRelease()
