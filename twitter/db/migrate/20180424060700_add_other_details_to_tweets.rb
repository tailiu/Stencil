class AddOtherDetailsToTweets < ActiveRecord::Migration[5.1]
  def change
    add_column :tweets, :tweet_media, :json
    add_column :tweets, :media_type, :string
  end
end
