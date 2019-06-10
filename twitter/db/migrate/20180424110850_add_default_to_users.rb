class AddDefaultToUsers < ActiveRecord::Migration[5.1]
  def change
    change_column_default :users, :bio, ""
    change_column_default :users, :protected, false
    change_column_default :users, :account_active, true
  end
end
