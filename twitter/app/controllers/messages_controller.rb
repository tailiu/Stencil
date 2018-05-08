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
            new_message = Message.new(content: content, user_id: user.id, conversation_id: conversation.id)
            new_message.save
        end
        
        render json: {result: result}
    end
end
