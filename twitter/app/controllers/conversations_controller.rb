class ConversationsController < ApplicationController
    def initialize
        @result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
    end

    def index 
        user = User.find(params[:id])

        conversation_participants = user.conversation_participants

        conversations = []

        for conversation_participant in conversation_participants do
            specific_conversation = conversation_participant.conversation
            specific_conversation_participants = specific_conversation.conversation_participants
            one_conversation = {
                "conversation" => specific_conversation,
                "conversation_participants" => specific_conversation_participants
            }
            conversations.push(one_conversation)
        end

        @result["conversations"] = conversations
        
        render json: {result: @result}
    end

    def new
        participants = params[:participants]

        # @conversation_participants = @user.conversation_participants

        conversation = Conversation.create()

        for participant in participants do
            user = User.find_by(handle: participant)
            conversation_participant = user.conversation_participants.create(conversation_id: conversation.id)
        end

        render json: {result: @result}
    end


end
