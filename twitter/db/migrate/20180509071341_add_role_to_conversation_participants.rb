class AddRoleToConversationParticipants < ActiveRecord::Migration[5.1]
  def change
    add_column :conversation_participants, :role, :string
  end
end
