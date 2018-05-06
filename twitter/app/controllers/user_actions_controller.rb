class UserActionsController < ApplicationController
    def index
        if params[:type] == "follow"
            @user = User.find(params[:id])
            @following_num = @user.user_actions.where(action_type: "follow").count
            @followed_num = User.joins("INNER JOIN user_actions ON user_actions.to_user_id = users.id").where("user_actions.action_type" => "follow").count

            @result = {
                # params: params,
                "followed_num" => @followed_num,
                "following_num" => @following_num
            }

            @result[:session] = session
            render json: {result: @result}
        end
    end

    def checkFollow
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        if params[:from_user_id].nil? || params[:to_user_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            result["success"] = true
            follow = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow")
            if follow.nil? || follow.empty?
                result["follow"] = false
            else
                result["follow"] = true
            end
        end
        
        render json: {result: result}
    end

    def checkBlock
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        if params[:from_user_id].nil? || params[:to_user_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            result["success"] = true
            block = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "block")
            if block.nil? || block.empty?
                result["block"] = false
            else
                result["block"] = true
            end
        end
        
        render json: {result: result}
    end

    def checkMute
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        if params[:from_user_id].nil? || params[:to_user_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            result["success"] = true
            mute = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "mute")
            if mute.nil? || mute.empty?
                result["mute"] = false
            else
                result["mute"] = true
            end
        end
        
        render json: {result: result}
    end

    def handleFollow
        result = {
            # params: params,
            "success" => true,
            "params" => params,
            "error" => {
            }
        }

        if params[:from_user_id].nil? || params[:to_user_id].nil? || params[:follow].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = User.where(id: params[:from_user_id]).first
            result["user"] = user
            if params[:follow] === "true"
                follow = UserAction.find_or_create_by(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow")
                # follow = user.build_user_action(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow")
                follow.save
                result["follow"] = follow
            else
                follow = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow").first
                UserAction.delete(follow.id)
                result["follow"] = false
            end
            result["success"] = true
        end
        render json: {result: result}
    end

    def handleBlock
        result = {
            # params: params,
            "success" => true,
            "params" => params,
            "error" => {
            }
        }

        if params[:from_user_id].nil? || params[:to_user_id].nil? || params[:block].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = User.where(id: params[:from_user_id]).first
            result["user"] = user
            if params[:block] === "true"
                block = UserAction.find_or_create_by(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "block")
                block.save
                result["block"] = block
            else
                block = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "block").first
                UserAction.delete(block.id)
                result["block"] = false
            end
            result["success"] = true
        end
        render json: {result: result}
    end
    
    def handleMute
        result = {
            # params: params,
            "success" => true,
            "params" => params,
            "error" => {
            }
        }

        if params[:from_user_id].nil? || params[:to_user_id].nil? || params[:mute].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = User.where(id: params[:from_user_id]).first
            result["user"] = user
            if params[:mute] === "true"
                mute = UserAction.find_or_create_by(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "mute")
                mute.save
                result["mute"] = mute
            else
                mute = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "mute").first
                UserAction.delete(mute.id)
                result["mute"] = false
            end
            result["success"] = true
        end
        render json: {result: result}
    end

end
