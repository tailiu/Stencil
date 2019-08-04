import proc_data as pd
import database as db
import graph as g
import datetime

def _log(fileName, l):
    f = open(fileName, "a")
    f.write("************************************\n")
    f.write(str(datetime.datetime.now()) + "\n")
    for data in l:
        f.write("[")
        f.write(", ".join([str(i) for i in data]))
        f.write("]\n")
    f.write("************************************\n")
    f.close()

def timeVsSize(stencilConnection, stencilCursor, destApp, startTime = None, endTime = None):
    if startTime == None or endTime == None:
        migrationIDs = pd.getAllMigrationIDs(stencilConnection, stencilCursor)
    else:
        migrationIDs = pd.getMigrationIDsBetweenTimestamps(stencilConnection, stencilCursor, startTime, endTime)
    time = []
    size = []
    for migrationID in migrationIDs:
        print migrationID
        l = pd.getLengthOfMigration(migrationID, stencilCursor)
        if  l == None:
            continue
        else:
            time.append(l)
            size.append(pd.getMigratedDataSize(destApp, migrationID, stencilCursor))
    _log("timeVsSize", [time, size])
    g.line(size, time, "Migration Size (KB)", "Migration Time (s)", "Migration Time Vs Migration Size")


def leftDataCumulativeGraph(stencilConnection, stencilCursor, destApp, srcApp, dataLeftInBrokenRows = False):
    migrationIDs = pd.getAllMigrationIDs(stencilConnection, stencilCursor)
    l = []
    for migrationID in migrationIDs:
        print migrationID
        time = pd.getLengthOfMigration(migrationID, stencilCursor)
        if time == None:
            continue
        leftData1 = pd.getSizeOfLeftDataInMigratedRows(srcApp, migrationID, stencilCursor)
        leftData2 = 0
        if not dataLeftInBrokenRows:
            leftData2 = pd.getSizeOfDataWithEntireRowLeft(srcApp, migrationID, stencilCursor)
        migratedData = pd.getMigratedDataSize(destApp, migrationID, stencilCursor)
        if leftData1 == None or leftData2 == None or migratedData == None:
            continue
        l.append((leftData1 + leftData2)/ float(migratedData + leftData1 + leftData2))
    if not dataLeftInBrokenRows:
        _log("allLeftData", [l])
        g.cumulativeGraph(l, "Percentage of Data Left", "Probability")
    else:
        _log("dataLeftInBrokenRows", [l])
        g.cumulativeGraph(l, "Percentage of Data Left in Broken Rows", "Probability")

def migratedDataInRowsBarGraph(l, step):
    l = pd.calComplement(l)
    x = range(0, 100, step)
    # print pd.getPercentageInIntervals(l, step/float(100))
    g.barGraph(x, pd.getPercentageInIntervals(l, step/float(100)), "Percentage of Data Migrated in Rows", "Frequency", step)

def leftDataBarGraph(l, step):
    x = range(0, 100, step)
    # print pd.getPercentageInIntervals(l, step/float(100))
    g.barGraph(x, pd.getPercentageInIntervals(l, step/float(100)), "Percentage of Data Left", "Frequency", step)

def anomaliesVsSize(stencilConnection, stencilCursor, srcApp):
    migrationIDs = pd.getAllMigrationIDs(stencilConnection, stencilCursor)
    print migrationIDs
    print pd.getMigratedDataSize(srcApp, migrationID, stencilCursor)

stencilDB, srcApp, destApp, migrationID = "stencil", "diaspora", "mastodon", 1017008071
stencilConnection, stencilCursor = db.connDB(stencilDB)
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-03-24 11:32:00", "2019-04-24 11:31:00") # 1 thread
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 11:32:00", "2019-04-24 12:09:00") # 5 threads
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 12:10:00", "2019-04-24 13:04:00") # 10 threads
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 13:09:00", "2019-04-24 15:45:00") # 20 threads
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 15:46:00", "2019-04-25 10:43:00") # 50 threads
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-25 10:43:00", "2019-09-24 15:45:00") # 100 threads
# leftDataCumulativeGraph(stencilConnection, stencilCursor, destApp, srcApp)
# leftDataCumulativeGraph(stencilConnection, stencilCursor, destApp, srcApp, True)

# g.allTimeVsSizeGraph()
# l = [0.474959160244, 0.304429522381, 0.362618944637, 0.701200316493, 0.131025582631, 0.305412530379, 0.370380350001, 0.323210377237, 0.351432578642, 0.540464020006, 0.356828861604, 0.907896500871, 0.496365782803, 0.136771134326, 0.372701746966, 0.229333333333, 0.355492554464, 0.426388888889, 0.222103463018, 0.204491296094, 0.206927134388, 0.289609644087, 0.202621518296, 0.194151352107, 0.613729104905, 0.130089899524, 0.585459815677, 0.319268653834, 0.196034381715, 0.302603004903, 0.219583816705, 0.620560538559, 0.484809171256, 0.226542459835, 0.194346105789, 0.357219035217, 0.586239064115, 0.201589678608, 0.73811706897, 0.461873343919, 0.348369207901, 0.111283109735, 0.33354833716, 0.310768151526, 0.463647796061, 0.200565184627, 0.685984160181, 0.377941126666, 0.258580691214, 0.38431525556, 0.430744777475, 0.600916002565, 0.195607788457, 0.156725526894, 0.307280976983, 0.198029505829, 0.17241220389, 0.413635296879, 0.357229154562, 0.454247300138, 0.504451543575, 0.180446995603, 0.253793449797, 0.657870417907, 0.289607723886, 0.318491624722, 0.338514910252, 0.384438542202, 0.163424124514, 0.186268008922, 0.320209842618, 0.433214143822, 0.294545454545, 0.371950385681, 0.355641873019, 0.328415337631, 0.199641986972, 0.203759398496, 0.509904299031, 0.46412558808, 0.0845097110374, 0.194484564001, 0.261344055141, 0.350568666759, 0.202172059417, 0.471030025626, 0.463326049215, 0.19786450358, 0.422825882958, 0.524160997474, 0.717921180917, 0.496770860243, 0.485549368907, 0.307117484542, 0.517786330767, 0.260677895612, 0.082013137558, 0.186745230733, 0.37737548352, 0.315062235926, 0.476489015207, 0.376740008981, 0.310210961385, 0.488481852387, 0.410911984238, 0.30653577282, 0.207414576671, 0.353575904095, 0.396303175024, 0.393961913609, 0.403064640523, 0.490277734774, 0.296833908764, 0.512586438725, 0.202037074761, 0.278523807272, 0.337497099544, 0.4797216204, 0.421501524151, 0.532551346733, 0.344421241899, 0.420342790712, 0.374725842977, 0.502886578482, 0.318336053816, 0.367562417614, 0.22679567711, 0.538272332023, 0.368786857624, 0.443158135669, 0.205388854868, 0.326944723079, 0.0855106888361, 0.469754442004, 0.472476549222, 0.165285678951, 0.437899159664, 0.153939302202, 0.58004643485, 0.259780132713, 0.433192038334, 0.406942487511, 0.318324834545, 0.539097219218, 0.665594939584, 0.170976510769, 0.173985857834, 0.222756559903, 0.34804199931, 0.361273038759, 0.532042311081, 0.61972400955, 0.322881481954, 0.631328543801, 0.146325016858, 0.135699373695, 0.352690224256, 0.342265488718, 0.326844130853, 0.362636300897, 0.201678468398, 0.200478878057, 0.497402612983, 0.683231189189, 0.302962573091, 0.177798010015, 0.410061888601, 0.177055320784, 0.196099621752, 0.31961772813, 0.418158567775, 0.724512156945, 0.133281188893, 0.349444941415, 0.255946859759, 0.298430933489, 0.480997654031, 0.394839918256, 0.398558690607, 0.683176052509, 0.521698551125, 0.474830429386, 0.184762220015, 0.488355031509, 0.316631236298, 0.368914033523, 0.164885600767, 0.202080593089, 0.360823524992, 0.680565926002, 0.257037913524, 0.332301790281, 0.586481624261, 0.324053901396, 0.340167462243, 0.64701754386, 0.368360009279, 0.118131868132, 0.410636702844, 0.327424270697, 0.161688869789, 0.429091430859, 0.447162783198, 0.172250859107, 0.226453517972, 0.371471698113, 0.334355050963, 0.380802297979, 0.295138512268, 0.580394334907, 0.30596071239, 0.338336630961, 0.198853574424, 0.409333022091, 0.492066720911, 0.281516003544, 0.238070069465, 0.289937652908, 0.201823435632, 0.379608358882, 0.345742316321, 0.383789514264, 0.141082656034, 0.574512707464, 0.158191030212, 0.830350118017, 0.392081039048, 0.318470855413, 0.195864828813, 0.195684821786, 0.401454784706, 0.34403175451, 0.307381978551, 0.223981533976, 0.195070316044, 0.372813352592, 0.307570066455, 0.194286997026, 0.411450898877, 0.445820004265, 0.540218576898, 0.330563026717, 0.488594177533, 0.0500974435401, 0.552145960887, 0.204432477216, 0.117398244214, 0.661230541142, 0.402610920762, 0.398950855881, 0.122467256241, 0.0971107544141, 0.235187540958, 0.868927852831, 0.677205507394, 0.318430022259, 0.410053097345, 0.281245011971, 0.296171334241, 0.630096682729, 0.176700822376, 0.172686230248, 0.412626275101, 0.104342475387, 0.27552, 0.320330226989, 0.2287694974, 0.374926542605, 0.227409361389, 0.21966618001, 0.166437069505, 0.320214936662, 0.556787867871, 0.33137641609, 0.278444401596, 0.697368421053, 0.477518077061, 0.378772990532, 0.4272702684, 0.195681310499, 0.275746105692, 0.419130881975, 0.194399843352, 0.579402011125, 0.195402783636, 0.321334293257, 0.310695478482, 0.488306993332, 0.179630159184, 0.182017934557, 0.26375728863, 0.637150328818, 0.676359923614, 0.362553549761, 0.573758452744, 0.354364234566, 0.409166119181, 0.69521117884, 0.498381133692, 0.301628358121, 0.137899087799, 0.471527000258, 0.180348140398, 0.22769121813, 0.364017190484, 0.369491074988, 0.488589765389, 0.665347991308, 0.218251859724, 0.440029305657, 0.505834339254, 0.117902350814, 0.485141424991, 0.465408857746, 0.33709208182, 0.385469299775, 0.121783741121, 0.569544551447, 0.223554603854, 0.178556797565, 0.343547039642, 0.166202531646, 0.455307951649, 0.422121304289, 0.338901762138, 0.602882551365]
# l = [0.223744486692, 0.19269327653, 0.230934829408, 0.208534926757, 0.0748403782429, 0.180297985076, 0.197291580032, 0.193165902419, 0.26095890411, 0.21697250702, 0.216575144805, 0.258899676375, 0.214759408469, 0.0746565006751, 0.203808359304, 0.229333333333, 0.194168354609, 0.225140712946, 0.136162939752, 0.144619126205, 0.128360352965, 0.20475540204, 0.120799710948, 0.162582144283, 0.253619909502, 0.0383865939205, 0.222295813435, 0.206730985641, 0.118170130705, 0.210703085033, 0.160642570281, 0.223089049679, 0.210794851951, 0.226542459835, 0.124370508686, 0.177613719604, 0.206674293817, 0.0999870146734, 0.203178224079, 0.23363233323, 0.20240719359, 0.0521005320125, 0.172640210115, 0.193428588182, 0.200765456018, 0.144717216915, 0.216538497729, 0.258426808503, 0.1720472964, 0.167308169853, 0.223117446545, 0.214322811472, 0.119199520148, 0.0620870667029, 0.213578538589, 0.112664802237, 0.134894252541, 0.199305648284, 0.210066984778, 0.224061087896, 0.216684856963, 0.0956808840375, 0.127213224805, 0.599855890477, 0.121551116334, 0.214860163127, 0.207981571563, 0.219192232207, 0.0966386554622, 0.0988049936578, 0.19312880525, 0.197166189363, 0.19151284491, 0.22492655791, 0.219896214549, 0.204402566288, 0.123884171565, 0.0894239036973, 0.216962622625, 0.208415742774, 0.0444971818451, 0.128947239982, 0.0875902992776, 0.210105740911, 0.129297678519, 0.182581191472, 0.217208080579, 0.150880574783, 0.18973828529, 0.218101333971, 0.208046350764, 0.248166730992, 0.208508660646, 0.210822852462, 0.176066161817, 0.122210152084, 0.0555555555556, 0.106753812636, 0.217386346978, 0.193528146289, 0.220474484122, 0.232300884956, 0.211492933496, 0.220691223587, 0.211244026421, 0.195401183662, 0.134699315478, 0.207779499307, 0.245412336945, 0.244557665586, 0.211356078335, 0.220497886999, 0.198831584057, 0.215975714319, 0.138925525367, 0.199260056152, 0.209605979515, 0.21948510499, 0.205026774918, 0.189393181045, 0.223172145959, 0.215455927536, 0.210420841683, 0.218168817415, 0.207995218501, 0.199147220202, 0.12480127186, 0.214405705697, 0.184767086025, 0.196970223411, 0.142403947832, 0.214352749797, 0.0855106888361, 0.22700814901, 0.220610525473, 0.102042173884, 0.19593701166, 0.080976652417, 0.22124334279, 0.134886951154, 0.217927527018, 0.23773460652, 0.19144646534, 0.219703866478, 0.212772521596, 0.0811186650185, 0.0870638239528, 0.131022602704, 0.203761371826, 0.213612298633, 0.226626076528, 0.1975371898, 0.207291899823, 0.220170785958, 0.0104957570344, 0.0678785857238, 0.250008129817, 0.211487088157, 0.214836151429, 0.233233150125, 0.0761760242792, 0.145843230404, 0.203792279866, 0.211794865652, 0.218285316924, 0.0889241190459, 0.213483699885, 0.134975518651, 0.127060329653, 0.187002968645, 0.198608872171, 0.214285714286, 0.0637093367131, 0.199155143156, 0.170851255377, 0.218038449713, 0.219994061905, 0.218164102403, 0.22917971196, 0.217550973226, 0.216855195388, 0.252909085434, 0.110836993802, 0.257821436423, 0.190927502605, 0.198527219605, 0.0802715956243, 0.0961668698361, 0.219754445274, 0.212960319726, 0.143851376303, 0.198680171885, 0.25127966496, 0.223128999351, 0.217955852347, 0.229709035222, 0.245392822502, 0.118131868132, 0.211778377824, 0.214130716918, 0.108941418294, 0.212061202368, 0.224686678818, 0.102229617304, 0.226453517972, 0.215841928381, 0.214633581823, 0.186535286564, 0.181066639782, 0.241084881969, 0.215735680705, 0.215261060425, 0.122594804267, 0.168190783575, 0.23802258163, 0.13043072907, 0.136118483007, 0.197767315786, 0.126397295381, 0.216864482606, 0.21146538083, 0.209965731985, 0.0633779264214, 0.20610849263, 0.105543001462, 0.238410596026, 0.120148436549, 0.215561557898, 0.133984580759, 0.128734884225, 0.213156503772, 0.2101597933, 0.179427044679, 0.138160557679, 0.127138562449, 0.216519647153, 0.232506004804, 0.106985134042, 0.211476262239, 0.225319361426, 0.194891768407, 0.210478456252, 0.220343186994, 0.0235682300259, 0.252766018365, 0.127481713689, 0.064778012685, 0.207972270364, 0.209045387725, 0.227192048278, 0.104728186386, 0.0637483355526, 0.142832454852, 0.208333333333, 0.221402214022, 0.198589101319, 0.227925507273, 0.215094997385, 0.200904233457, 0.236910096655, 0.139936601277, 0.118738192019, 0.215679160964, 0.0417567948839, 0.27552, 0.211088974197, 0.2287694974, 0.232009626955, 0.0858013937282, 0.150515463918, 0.100540540541, 0.218290912344, 0.233294012034, 0.179559332527, 0.200762305316, 0.101372756072, 0.20705060515, 0.220261760353, 0.209755250825, 0.144463804847, 0.196417322199, 0.231371214699, 0.142624932272, 0.195933740073, 0.144797597167, 0.192106535402, 0.184396485867, 0.218753030156, 0.126731266805, 0.107323395982, 0.191528102927, 0.219763021032, 0.200954168184, 0.233914647938, 0.214824514886, 0.223644403072, 0.213287109786, 0.180561806677, 0.221600762202, 0.237153452234, 0.0759300282637, 0.210638708432, 0.124377603902, 0.13835335019, 0.18192478432, 0.24301242236, 0.221092634863, 0.204076367389, 0.0700063211125, 0.223535019427, 0.222483494619, 0.079275198188, 0.226049515608, 0.224009343861, 0.171696677657, 0.198032580728, 0.0893690154677, 0.231365550717, 0.137241838774, 0.105075235909, 0.212557783692, 0.100996314999, 0.27271396417, 0.209698934293, 0.248583591094, 0.579454427365]
# migratedDataInRowsBarGraph(l, 10)

# print len(pd.getMigrationIDsBetweenTimestamps(stencilConnection, stencilCursor, "2019-04-24 11:32:00", "2019-04-24 12:09:00"))
# l = {0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.1, 0.1, 0.1, 0.1, 0.2, 0.2, 0.8}
# g.cumulativeGraph(l, "Percentage of Data Left", "Probability")
# print pd.getSizeOfLeftDataInMigratedRows(srcApp, migrationID, stencilCursor)
# print pd.getSizeOfDataWithEntireRowLeft(srcApp, migrationID, stencilCursor)
# print getAllMigrationIDs(stencilConnection, stencilCursor)
# print getMigratedDataSize(destApp, migrationID, stencilCursor)
# print getLengthOfMigration(migrationID, stencilCursor)


anomaliesVsSize()

db.closeDB(stencilConnection)

