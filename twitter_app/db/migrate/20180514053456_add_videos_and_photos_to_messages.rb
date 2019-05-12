class AddVideosAndPhotosToMessages < ActiveRecord::Migration[5.1]
  def change
    add_column :messages, :message_media, :string
  end
end
