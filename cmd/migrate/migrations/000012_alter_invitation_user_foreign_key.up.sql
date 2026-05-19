ALTER TABLE user_invitations
ADD CONSTRAINT user_invitations_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;