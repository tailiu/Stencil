class RenameTypeToActionType < ActiveRecord::Migration[5.1]
  def change
    rename_column :user_actions, :type, :action_type
  end
end
