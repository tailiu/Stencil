class ConversationsController < ApplicationController
    def index 
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

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

        result["conversations"] = conversations
        
        render json: {result: result}
    end

    def new
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        participants = params[:participants]

        userArr = []
        error = false

        for participant in participants do 
            user = User.find_by(handle: participant)
            if user == nil
                result['success'] = false
                result['error'] = 'Invalid handle'
                break
            end
            userArr.push(user)
        end

        if result['success']
            conversation = Conversation.create()
            for user in userArr do
                conversation_participant = user.conversation_participants.create(conversation_id: conversation.id)
            end
        end

        render json: {result: result}
    end


end
