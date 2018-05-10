class ConversationsController < ApplicationController    
    def index 
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        user = User.find_by_id(params[:id])

        conversation_participants = user.conversation_participants.order(created_at: :desc)

        conversations = []

        for conversation_participant in conversation_participants do
            one_conversation = conversation_participant.conversation
            one_conversation_participants = one_conversation.conversation_participants

            conversation = {
                "conversation" => one_conversation,
                "conversation_participants" => [],
                "conversation_state" => 'normal'
            }
            for one_conversation_participant in one_conversation_participants do
                user = one_conversation_participant.user
                conversation["conversation_participants"].push(user)
            end
            if one_conversation.conversation_type == "not_group" 
                block_one = UserAction.where(from_user_id: conversation["conversation_participants"][0], 
                    to_user_id: conversation["conversation_participants"][1], action_type: "block")
                block_two = UserAction.where(from_user_id: conversation["conversation_participants"][1], 
                    to_user_id: conversation["conversation_participants"][0], action_type: "block")
                if (block_one.length != 0 || block_two.length != 0)
                    conversation['conversation_state'] = 'blocked'
                end
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
            },
            "conversation_state" => ""            
        }

        parameters = JSON.parse(params[:participants])
        participants = parameters['participants']

        userArr = []
        error = false

        if participants.length == 0
            result['success'] = false
            result['error'] = 'No handle submitted'
        else 
            for participant in participants do 
                user = User.find_by(handle: participant)
                if user == nil
                    result['success'] = false
                    result['error'] = 'Invalid handle'
                    break
                end
                userArr.push(user)
            end
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
                        next 
                    end

                    if conversation.conversation_participants.size == 1
                        existed = true
                        result["conversation"] = conversation
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
                            block_one = UserAction.where(from_user_id: conversation_participants[0], 
                                to_user_id: conversation_participants[1], action_type: "block")
                            block_two = UserAction.where(from_user_id: conversation_participants[1], 
                                to_user_id: conversation_participants[0], action_type: "block")
                            if (block_one.length != 0 || block_two.length != 0)
                                result["conversation_state"] = "blocked"
                            end
                            result["conversation"] = conversation
                        end
                    end
                end 
            end

            if !existed
                conversation_creator = parameters['conversation_creator']
                creator = User.find_by(handle: conversation_creator)

                conversation = Conversation.create(conversation_type: conversation_type)
                for user in userArr do
                    if user.id == creator.id 
                        user.conversation_participants.create(conversation_id: conversation.id, role: 'creator')
                    else
                        user.conversation_participants.create(conversation_id: conversation.id, role: 'normal')
                    end
                end

                result["conversation"] = conversation
            end
            
        end

        render json: {result: result}
    end

    
    def leaveConversation
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        conversation = Conversation.find_by_id(params[:conversation_id])
        user = User.find_by_id(params[:user_id])

        if user == nil || conversation == nil || !conversation.conversation_participants.find_by(user_id: user.id)
            result['success'] = false
            result['error'] = 'No such conversation or user, or this user does not belong to this conversation'
        end

        if result['success']
            conversation_participants = conversation.conversation_participants
            conversation_participant = conversation_participants.find_by(user_id: user.id)
            conversation_participants.delete(conversation_participant)

            if conversation_participants.empty?
                conversation.messages.clear()
                conversation.delete()
            end
        end
        
        render json: {result: result}
    end

    def blockInGroupConversation
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        from_user = User.find_by_id(params[:from_user_id])
        to_user = User.find_by_id(params[:to_user_id])

        from_user_conversation_participants = from_user.conversation_participants

        for from_user_conversation_participant in from_user_conversation_participants do 
            conversation = from_user_conversation_participant.conversation
            if conversation.conversation_type == 'not_group'
                next
            end
            conversation_participants = conversation.conversation_participants
            to_user_conversation_participant = conversation_participants.find_by(user_id: to_user.id)
            if to_user_conversation_participant == nil
                next 
            end
            if to_user_conversation_participant.role == 'creator'
                from_user_conversation_participant.delete()
                next
            end
            to_user_conversation_participant.delete()
        end

        render json: {result: result}

    end
end
