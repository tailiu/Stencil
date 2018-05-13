class MessagesController < ApplicationController
    def index
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        conversation = Conversation.find_by(id: params[:id])
        messages = nil
        if conversation == nil
            result["success"] = false
            result["error"] = "No such conversation"
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
        user = User.find_by_id(params[:user_id])
        conversation = Conversation.find_by_id(params[:conversation_id])
        content = params[:content]

        if user == nil || conversation == nil
            result["success"] = false
            result["error"] = "No such conversation or user"
        end

        if result["success"]
            conversation_participants = conversation.conversation_participants
            userInConversation = false
            for conversation_participant in conversation_participants do
                if conversation_participant.user_id == user.id
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
            if conversation.conversation_type == 'not_group'
                block_one = UserAction.find_by(from_user_id: conversation_participants[0].user_id, 
                    to_user_id: conversation_participants[1].user_id, action_type: "block")
                block_two = UserAction.find_by(from_user_id: conversation_participants[1].user_id, 
                    to_user_id: conversation_participants[0].user_id, action_type: "block")
                puts '*******************************'
                puts block_one != nil
                puts block_two != nil
                puts '*******************************'
                if (block_one != nil || block_two != nil)
                    result["success"] = false
                    result["error"] = "You have been blocked in this conversation"
                end
            end
        end

        if result["success"]
            new_message = Message.new(content: content, user_id: user.id, conversation_id: conversation.id)
            new_message.save
        end

        render json: {result: result}
    end
end
