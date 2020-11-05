- DoTransactions example
```
var pid string

// the transaction unit function
testFunc1 := func(db *gorm.DB) error {
    player := Player{
        Name: "test name",
    }
    err := db.Create(&player).Error
    if err != nil {
        // You can change err before return
        return errors.New(fmt.Sprintf("Create Player error: %s", err.Error()))
    }
    pid = player.ID

    return nil
}

// the transaction unit function
testFunc2 := func(db *gorm.DB) error {
    sp := Sport{
        Name:     "test sport name",
        PlayerId: pid, // pid is uesed here.
    }

    err = db.Create(&sp).Error
    if err != nil {
        return err
    }
    return nil
}

// do the transactions
if err := kipleDb.DoTransactions(testFunc1, testFunc2); err != nil {
    // handle the error
}
```