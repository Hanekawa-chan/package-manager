package client

func createDirCmd(name, ver string) string {
	return "mkdir -p " + name + "/" + ver
}

func sendFileCmd(name, ver string) string {
	return "cat - > " + name + "/" + ver + "/data.zip"
}

func getPackageCmd(name, ver string) string {
	return "cat " + name + "/" + ver + "/data.zip"
}

func findFilesCmd(name string) string {
	return "find ./" + name + " -maxdepth 1 -name '*.*'"
}
