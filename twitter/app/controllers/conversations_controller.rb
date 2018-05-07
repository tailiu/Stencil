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
            existed = false
            conversation_type = 'group'

            if userArr.length == 1
                conversation_type = 'not_group'
                conversation_participants = userArr[0].conversation_participants
                for conversation_participant in conversation_participants do
                    conversation = conversation_participant.conversation
                    if conversation.conversation_type != conversation_type
                        continue 
                    end

                    if conversation.conversation_participants.size == 1
                        existed = true
                        result["conversationID"] = conversation.id
                        break
                    end
                end
            end

            if userArr.length == 2
                conversation_type = 'not_group'
                user_one_conversation_participants = userArr[0].conversation_participants
                for user_one_conversation_participant in user_one_conversation_participants do
                    conversation = user_one_conversation_participant.conversation
                    if conversation.conversation_type != conversation_type
                        next 
                    end

                    conversation_participants = conversation.conversation_participants
                    if conversation_participants.size == 2 
                        if conversation_participants.exists?(user_id: userArr[1].id) 
                            existed = true
                            result["conversationID"] = conversation.id
                        end
                    end
                end 
            end

            if !existed
                conversation = Conversation.create(conversation_type: conversation_type)
                for user in userArr do
                    conversation_participant = user.conversation_participants.create(conversation_id: conversation.id)
                end
                result["conversationID"] = conversation.id
            end
            
        end

        render json: {result: result}
    end

end