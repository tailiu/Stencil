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
        if (conversation!=nil) then
            messages = Message.joins("INNER JOIN users ON messages.user_id = users.id")
                        .where('messages.conversation_id' => conversation.id)
                        .select("users.name, users.avatar, messages.*")
        end
        result["messages"] = messages

        render json: {result: result}
    end

    def new
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        new_message = Message.new(content: params[:content], user_id: params[:user_id], conversation_id: params[:conversation_id])
        new_message.save
        
        render json: {result: result}
    end
end
