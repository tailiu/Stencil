class DropTweetMediaTable < ActiveRecord::Migration[5.1]
  def change
    drop_table :tweet_media
  end
end
