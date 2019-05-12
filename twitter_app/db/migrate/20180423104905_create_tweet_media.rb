class CreateTweetMedia < ActiveRecord::Migration[5.1]
  def change
    create_table :tweet_media do |t|
      t.references :tweet, foreign_key: true
      t.string :media_type

      t.timestamps
    end
  end
end
