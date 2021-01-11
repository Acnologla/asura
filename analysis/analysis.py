import json
import matplotlib as mpl
import matplotlib.patches as mpatches
import matplotlib.pyplot as plt
import numpy as np
from matplotlib.ticker import FormatStrFormatter

with open('./data.json') as f:
  data = json.load(f)

mpl.style.use("seaborn")

lines = [
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
    [[],[]],
]

for x in data:
    lines[x["Class"][1]-1][0].append(x["Level"][1])
    lines[x["Class"][1]-1][1].append(x["Wins"][0]/100)


plt.legend(handles=[
  mpatches.Patch(color='red', label='NORMAL'),
  mpatches.Patch(color='g', label='PAPEL'),
  mpatches.Patch(color='b', label='PEDRA'),
  mpatches.Patch(color='y', label='CORTANTE'),
  mpatches.Patch(color='pink', label='FOGO'),
  mpatches.Patch(color='cyan', label='GELO'),
  mpatches.Patch(color='black', label='MAGMA'),
  mpatches.Patch(color='orange', label='AGUA'),
  mpatches.Patch(color='brown', label='NATUREZA'),
  mpatches.Patch(color='grey', label='LUZ'),
  mpatches.Patch(color='#85FF28', label='ACIDO'),
  mpatches.Patch(color='#7024F2', label='ELETRICO'),
  mpatches.Patch(color='#A06D9E', label='GRAVIDADE')
])

plt.plot(lines[0][0], lines[0][1], 'r',marker='o')
plt.plot(lines[1][0], lines[1][1], 'g',marker='o')
plt.plot(lines[2][0], lines[2][1], 'b',marker='o')
plt.plot(lines[3][0], lines[3][1], 'y',marker='o')
plt.plot(lines[4][0], lines[4][1], 'pink',marker='o')
plt.plot(lines[5][0], lines[5][1], 'cyan',marker='o')
plt.plot(lines[6][0], lines[6][1], 'black',marker='o')
plt.plot(lines[7][0], lines[7][1], 'orange',marker='o')
plt.plot(lines[8][0], lines[8][1], 'brown',marker='o')
plt.plot(lines[9][0], lines[9][1], 'grey',marker='o')
plt.plot(lines[15][0], lines[15][1], '#85FF28',marker='o')
plt.plot(lines[16][0], lines[16][1], '#7024F2',marker='o')
plt.plot(lines[17][0], lines[17][1], '#A06D9E',marker='o')

plt.plot([0,30],[50,50], 'm', linestyle='--')
ax=plt.gca()

ax.yaxis.set_major_locator(mpl.ticker.MultipleLocator(10))
ax.xaxis.set_major_locator(mpl.ticker.MultipleLocator(2))


plt.ylabel('Percentage')
plt.show()