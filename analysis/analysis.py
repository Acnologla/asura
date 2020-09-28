import json
import matplotlib.pyplot as plt

with open('./data.json') as f:
  data = json.load(f)

lines = [
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
]

for x in data:
    lines[x["Class"][1]-1][0].append(x["Level"][1])
    lines[x["Class"][1]-1][1].append(x["Wins"][0]/100)

plt.plot(lines[0][0], lines[0][1], 'r')
#plt.plot(lines[3][0], lines[3][1], 'p')

plt.ylabel('Percentage')
plt.show()