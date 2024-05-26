package handlers

// func Regenerate(config config.Config, storageProvider providers.StorageProvider) WannabeHandler {
// 	return func(ctx *fiber.Ctx) error {
// 		if !config.StorageProvider.Regenerate {
// 			return internalError(ctx, fmt.Errorf("regenerate set to false in config"))
// 		}

// 		regenCount := 0
// 		regenHashes := []string{}
// 		failedCount := 0
// 		failedHashes := []string{}

// 		hashes, err := storageProvider.GetHashes()
// 		if err != nil {
// 			return internalError(ctx, err)
// 		}

// 		// REVIEW mem issue in case of too many records ?
// 		records, err := storageProvider.ReadRecords(hashes)
// 		if err != nil {
// 			return internalError(ctx, err)
// 		}

// 		for _, encodedRecord := range records {
// 			var record recordEntities.Record

// 			oldHash := record.Request.Hash

// 			err := json.Unmarshal(encodedRecord, &record)
// 			if err != nil {
// 				failedCount++
// 				failedHashes = append(failedHashes, oldHash)

// 				continue
// 			}

// 			requestBody, err := json.Marshal(record.Request.Body)
// 			if err != nil {
// 				failedCount++
// 				failedHashes = append(failedHashes, oldHash)

// 				continue
// 			}

// 			curlPayload := curlEntities.GenerateCurlPayload{
// 				HttpMethod:     record.Request.HttpMethod,
// 				Path:           record.Request.Path,
// 				Query:          record.Request.Query,
// 				RequestHeaders: record.Request.Headers,
// 				RequestBody:    requestBody,
// 			}

// 			curl, err := curl.GenerateCurl(config, curlPayload)
// 			if err != nil {
// 				failedCount++
// 				failedHashes = append(failedHashes, oldHash)

// 				continue
// 			}

// 			hash, err := hash.GenerateHash(curl)
// 			if err != nil {
// 				failedCount++
// 				failedHashes = append(failedHashes, oldHash)

// 				continue
// 			}

// 			isDuplicateHash := checkDuplicates(hashes, hash)
// 			isDuplicateRegenHash := checkDuplicates(regenHashes, hash)
// 			if isDuplicateHash || isDuplicateRegenHash {
// 				continue
// 			}

// 			record.Request.Hash = hash
// 			record.Request.Curl = curl
// 			record.Metadata.RegeneratedAt = recordEntities.Timestamp{
// 				Unix: time.Now().Unix(),
// 				UTC:  time.Now().UTC(),
// 			}

// 			encodedRecordRegen, err := json.Marshal(record)
// 			if err != nil {
// 				failedCount++
// 				failedHashes = append(failedHashes, oldHash)

// 				continue
// 			}

// 			err = storageProvider.InsertRecords([]string{hash}, [][]byte{encodedRecordRegen})
// 			if err != nil {
// 				failedCount++
// 				failedHashes = append(failedHashes, oldHash)

// 				continue
// 			}

// 			regenCount++
// 			regenHashes = append(regenHashes, hash)
// 		}

// 		return ctx.Status(fiber.StatusCreated).JSON(RegenerateResponse{
// 			Message:           fmt.Sprintf("%v records succeeded in regenerating, %v records failed in regenerating", regenCount, failedCount),
// 			RegeneratedHashes: regenHashes,
// 			FailedHashes:      failedHashes,
// 		})
// 	}
// }
