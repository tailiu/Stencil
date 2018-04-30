class AddForeignKeyToActions < ActiveRecord::Migration[5.1]
  def change
    add_foreign_key :user_actions, :users, column: :from_user_id
  end
end
