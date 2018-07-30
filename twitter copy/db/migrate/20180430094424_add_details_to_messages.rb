class AddDetailsToMessages < ActiveRecord::Migration[5.1]
  def change
    add_column :messages, :content, :text 
    add_column :messages, :media_type, :string
    add_reference :messages, :conversation, foreign_key: true
    add_reference :messages, :user, foreign_key: true
  end
end
