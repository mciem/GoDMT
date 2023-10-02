with open("tokens.txt", "r") as f:
    lines = f.read().splitlines()

with open("tokens.txt", "w") as f:
    for line in lines:
        f.write(line.split(":")[2]+"\n")