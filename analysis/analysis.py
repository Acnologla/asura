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
]

for x in data:
    lines[x["Class"][1]-1][0].append(x["Level"][1])
    lines[x["Class"][1]-1][1].append(x["Wins"][0]/100)


plt.legend(handles=[
  mpatches.Patch(color='red', label='NORMAL'),
  mpatches.Patch(color='g', label='PAPEL'),
  mpatches.Patch(color='b', label='PEDRA'),
  mpatches.Patch(color='y', label='CORTANTE'),
])

plt.plot(lines[0][0], lines[0][1], 'r',marker='o')
plt.plot(lines[1][0], lines[1][1], 'g',marker='o')
plt.plot(lines[2][0], lines[2][1], 'b',marker='o')
plt.plot(lines[3][0], lines[3][1], 'y',marker='o')
plt.plot([0,50],[50,50], 'm', linestyle='--')
ax=plt.gca()

ax.yaxis.set_major_locator(mpl.ticker.MultipleLocator(10))
ax.xaxis.set_major_locator(mpl.ticker.MultipleLocator(2))


plt.ylabel('Percentage')
plt.show()