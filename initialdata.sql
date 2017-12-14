INSERT INTO login_user
(
  username,
  password,
  name,
  roles,
  is_enabled,
  created,
  last_login,
  is_sync
)
VALUES
(
  'administrator',
  '$2a$04$mGSEn9GYFUlQrerysW0Vr.KbFxSVc2qpMKpFF4h6jzFyaBU0vgsVm',
  'Administrator',
  'ADMIN',
  TRUE,
  TIMESTAMP '2017-12-13 19:31:47.820',
  TIMESTAMP '2017-12-13 19:31:47.820',
  FALSE
);

