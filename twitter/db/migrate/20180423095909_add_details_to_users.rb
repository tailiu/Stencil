class AddDetailsToUsers < ActiveRecord::Migration[5.1]
  def change
    add_column :users, :name, :string
    add_column :users, :handle, :text
    add_column :users, :bio, :text
    add_column :users, :protected, :boolean
    add_column :users, :account_active, :boolean
  end
end
