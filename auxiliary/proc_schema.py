import re

schemaFile = open('/home/nyuad/mastodon.sql', 'r')
output = open('/home/nyuad/proc_mastodon.sql', 'w')

stack = []
tables = {}
copy = ''
constraint = ''
index = ''
result = ''

for line in schemaFile:
    if len(stack):
        catogory, key = stack[0]
        if catogory == 'tables':
            tables[key] += line
            if line.find(');') != -1:
                stack.pop(0)
        if catogory == 'copy':
            copy += line
            if line.find('\.') != -1:
                stack.pop(0)
        if catogory == 'constraint':
            if line.find('PRIMARY KEY') != -1:
                primaryKey = line[line.find('(') + 1:line.find(')')]
                createTable = tables[key]
                createTablelines = createTable.split('\n')
                for i, createTableline in enumerate(createTablelines):
                    if re.search(r'\b%s\b'%(primaryKey), createTableline):
                        createTableline = createTableline.replace('bigint', 'serial8')
                        if createTableline.find(',') != -1:
                            createTablelines[i] = createTableline[:-1] + ' PRIMARY KEY,'
                        else:
                            createTablelines[i] = createTableline + ' PRIMARY KEY'
                tables[key] = '\n'.join(createTablelines)
            elif line.find('FOREIGN KEY') != -1:
                constraint += 'ALTER TABLE ONLY ' + key + line
            stack.pop(0)

    elif line.find('CREATE TABLE') != -1:
        preLen = len('CREATE TABLE')
        lineLen = len(line)
        tableName = line[preLen + 1:lineLen - 3]
        tables[tableName] = line
        stack.append(['tables', tableName])

    elif line.find('COPY') != -1:
        copy += line
        stack.append(['copy', None])
    
    elif line.find('ALTER TABLE') != -1 and line.find('OWNER TO') == -1 and line.find('SET') == -1:
        stack.append(['constraint', line[len('ALTER TABLE ONLY '):-1]])

    elif line.find('CREATE INDEX') != -1 or line.find('CREATE UNIQUE INDEX') != -1:
        index += line


for key in tables:
    result += tables[key]

result += copy
result += index
result += constraint
output.write(result)

schemaFile.close()
output.close()
