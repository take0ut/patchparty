package db

const GetPatchById string = `
	SELECT *
	FROM patches
	JOIN users ON users.id = patches.user_id
	WHERE patches.id = $1
`

const CreatePatch string = `
	INSERT INTO patches (name, description, device_id, user_id, downloadeds, source_url) 
	VALUES ($1, $2, $3, $4, $5, $6)
`
