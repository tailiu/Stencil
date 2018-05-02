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
            one_conversation = conversation_participant.conversation
            one_conversation_participants = one_conversation.conversation_participants

            conversation = {
                "conversation" => one_conversation,
                "conversation_participants" => []
            }
            for one_conversation_participant in one_conversation_participants do
                user = one_conversation_participant.user
                conversation["conversation_participants"].push(user)
            end

            conversations.push(conversation)
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
