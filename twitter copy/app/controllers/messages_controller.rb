class MessagesController < ApplicationController
    def index
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        user = User.find_by(id: params[:user_id])
        conversation = Conversation.find_by(id: params[:conversation_id])
        messages = nil

        if conversation == nil || user == nil
            result["success"] = false
            result["error"] = "No such conversation or user"
        end

        if result["success"]
            exist = false
            conversation_participants = conversation.conversation_participants
            for conversation_participant in conversation_participants do
                if conversation_participant.user_id == user.id 
                    exist = true
                end
            end
            if !exist
                result["success"] = false
                result["error"] = "The user is not in the conversation"
            end
        end

        if result["success"]
            messages = Message.joins("INNER JOIN users ON messages.user_id = users.id")
                        .where('messages.conversation_id' => conversation.id)
                        .select("users.name, users.avatar, messages.*")
                        .order(created_at: :asc)
            result["messages"] = messages
        end

        render json: {result: result}
    end

    def new
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        message_media = params[:media]
        media_type = params[:media_type]

        user = User.find_by(id: params[:user_id])
        conversation = Conversation.find_by(id: params[:conversation_id])
        content = params[:content]

        if user == nil || conversation == nil
            result["success"] = false
            result["error"] = "No such conversation or user"
        end

        if media_type != 'video' && media_type != 'photo' && media_type != 'null'
            result["success"] = false
            result["error"] = "Wrong media type"
        end

        if result["success"]
            conversation_participants = conversation.conversation_participants
            userInConversation = false
            for conversation_participant in conversation_participants do
                if conversation_participant.user_id == user.id
                    participant = conversation_participant
                    userInConversation = true
                    break
                end
            end
            if !userInConversation
                result["success"] = false
                result["error"] = "This user is not in the conversation"
            end
        end

        if result["success"]
            if conversation.conversation_type == 'not_group' && conversation_participants.length == 2
                block_one = UserAction.find_by(from_user_id: conversation_participants[0].user_id, 
                    to_user_id: conversation_participants[1].user_id, action_type: "block")
                block_two = UserAction.find_by(from_user_id: conversation_participants[1].user_id, 
                    to_user_id: conversation_participants[0].user_id, action_type: "block")
                if (block_one != nil || block_two != nil)
                    result["success"] = false
                    result["error"] = "You have been blocked in this conversation"
                end
            end
        end

        if result["success"]
            new_message = Message.new(
                                content: content, 
                                user_id: user.id, 
                                conversation_id: conversation.id,
                                media_type: media_type,
                                message_media: message_media
                            )
            new_message.save
            participant.saw_messages_until = new_message.created_at
            participant.save
            result["newMessage"] = new_message
        end

        render json: {result: result}
    end

    def setSawMessagesUntil
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        user = User.find_by(id: params[:user_id]) 
        conversation = Conversation.find_by(id: params[:conversation_id])
        message = Message.find_by(id: params[:message_id])
        
        if conversation == nil || user == nil || message == nil
            result["success"] = false
            result["error"] = "No such conversation, user or message"
        end

        if result["success"]
            conversation_participants = conversation.conversation_participants
            for conversation_participant in conversation_participants do
                if conversation_participant.user_id == user.id
                    conversation_participant.saw_messages_until = message.created_at
                    conversation_participant.save
                end
            end
        end

        render json: {result: result}
    end

end
