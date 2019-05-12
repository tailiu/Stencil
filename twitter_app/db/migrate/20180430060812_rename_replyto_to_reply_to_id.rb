class RenameReplytoToReplyToId < ActiveRecord::Migration[5.1]
  def change
    rename_column :tweets, :reply_to, :reply_to_id
  end
end
