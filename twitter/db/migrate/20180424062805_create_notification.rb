class CreateNotification < ActiveRecord::Migration[5.1]
  def change
    create_table :notifications do |t|
      t.string :notification_type
      t.references :user, foreign_key: true
      t.bigint :from_user
      t.bigint :tweet
      t.boolean :is_seen
    end
  end
end
