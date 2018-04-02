import subprocess, json, glob

nodes = ["00", "01", "02"]
addrs = []  
for node in nodes:
    mkdircmd = "mkdir -p data/" + node
    subprocess.run(mkdircmd, shell=True)
    print("creating account " + node)
    createcmd = "./build/bin/geth account new --datadir data/" + node + " --password ./pass"
    subprocess.run(createcmd, shell=True)
    file = "data/" + node + "/keystore/*"
    for f in glob.glob(file):
        addr = json.load(open(f))['address']
        addrs.append(addr)

print("Creating genesis.json")
genesis = json.load(open("genesis.json.template"))
genesis['config']['thw']['accounts']= addrs 

with open('genesis.json', 'w') as fp:
    json.dump(genesis, fp, indent=4)
