import matplotlib.pyplot as plt
import numpy as np
import proc_data as pd
import database as db

def lengthVsSize(stencilConnection, stencilCursor, destApp):
    migrationIDs = pd.getAllMigrationIDs(stencilConnection, stencilCursor)
    length = []
    size = []
    for migrationID in migrationIDs:
        print migrationID
        l = pd.getLengthOfMigration(migrationID, stencilCursor)
        if  l == None:
            continue
        else:
            length.append(l)
            size.append(pd.getMigratedDataSize(destApp, migrationID, stencilCursor))
    order = np.argsort(size)
    xs = np.array(size)[order]
    ys = np.array(length)[order]

    ax = plt.axes()
    ax.plot(xs, ys)
    plt.xlabel("Migration Size (bytes)")
    plt.ylabel("Migration Length (seconds)")
    plt.show()


stencilDB, srcApp, destApp, migrationID = "stencil", "diaspora", "mastodon", 2108391555
stencilConnection, stencilCursor = db.connDB(stencilDB)
lengthVsSize(stencilConnection, stencilCursor, destApp)

# print getAllMigrationIDs(stencilConnection, stencilCursor)
# print getMigratedDataSize(destApp, migrationID, stencilCursor)
# print getLengthOfMigration(migrationID, stencilCursor)
db.closeDB(stencilConnection)

