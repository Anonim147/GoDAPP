package temp

func method3(db *gorm.DB) {
    records := make([]Record, 0)
    defer timeTrack(time.Now(), "Method 3")
    count := 1000000
    bucketSize := 100000
    resultCount := 0
    resultChannel := make(chan []Record, 0)
    for beginID :=1; beginID <=count; beginID += bucketSize {
        endId := beginID +bucketSize
        go func(beginId int, endId int) {
            currentRecords := make([]Record, 0)
            db.Table("records").Select("id").Where("id>=? and id<?", beginId, endId).Find(&currentRecords)
            resultChannel <- currentRecords
        }(beginID, endId)
        resultCount += 1
    }
    for i:=0; i<resultCount; i++{
        currentRecords := <- resultChannel
        records = append(records, currentRecords...)
    }
}