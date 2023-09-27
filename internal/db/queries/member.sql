-- name: CreateMember :exec
INSERT INTO member (merchant_id, email, password, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: GetMemberInfo :one
SELECT real_name,
       email,
       mobile,
       sex,
       birthday,
       city,
       district,
       address,
       zip_code
FROM member
WHERE id = ?;

-- name: UpdateMemberPassword :exec
UPDATE member
SET password   = ?,
    updated_at = ?
WHERE id = ?
  AND password = sqlc.arg('old_password');

-- name: MemberLogin :one
SELECT *
FROM member
WHERE merchant_id = ?
  AND email = ?
  AND password = ?;

-- name: ResetPassword :exec
UPDATE member
SET password = ?
WHERE merchant_id = ?
  AND email = ?;

-- name: CheckEmailExist :one
SELECT 1
FROM member
WHERE merchant_id = ?
  AND email = ?;