class ConversationsController < ApplicationController
    def index 
        @user = User.find(params[:id])

        @conversation_participants = @user.conversation_participants

        for conversation_participant in @conversation_participants do
            @conversation = conversation_participant.conversation
            puts @conversation.id
        end

        @result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        
        render json: {result: @result}
    end

    def new
        @user = User.find(params[:id])
        @participants = params[:participants]

        # @conversation_participants = @user.conversation_participants

        @conversation = Conversation.create()

        for participant in @participants do
            @user = User.find_by(handle: participant)
            @conversation_participant = @user.conversation_participants.create(conversation_id: @conversation.id)
        end

        @result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        render json: {result: @result}
    end


end
