class AddDetailsToTweets < ActiveRecord::Migration[5.1]
  def change
    add_column :tweets, :type, :string
    add_column :tweets, :content, :text
    add_column :tweets, :reply_to, :integer
    add_reference :tweets, :user, foreign_key: true
  end
end
