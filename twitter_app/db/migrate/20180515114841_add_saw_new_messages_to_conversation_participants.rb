class AddSawNewMessagesToConversationParticipants < ActiveRecord::Migration[5.1]
  def change
    add_column :conversation_participants, :saw_new_messages, :boolean
  end
end
