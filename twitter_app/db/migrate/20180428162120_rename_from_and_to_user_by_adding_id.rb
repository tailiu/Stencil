class RenameFromAndToUserByAddingId < ActiveRecord::Migration[5.1]
  def change
    rename_column :user_actions, :from_user, :from_user_id
    rename_column :user_actions, :to_user, :to_user_id
  end
end
