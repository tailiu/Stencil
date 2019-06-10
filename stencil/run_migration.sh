for i in $(seq 4851 4900)
do
    go run migration_handler.go $i
    sleep 5
done