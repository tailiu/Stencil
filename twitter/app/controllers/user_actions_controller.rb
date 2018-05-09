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
                follow = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow_pending")
                if follow.nil? || follow.empty?
                    result["follow"] = false
                else
                    result["follow"] = "pending"
                end
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

    def checkTwoWayBlock
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
            is_blocked = UserAction.where(from_user_id: params[:to_user_id], to_user_id: params[:from_user_id], action_type: "block")
            if (block.nil? || block.empty?) && (is_blocked.nil? || is_blocked.empty?)
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
    
    def getFollowRequests
        result = {
            # params: params,
            "success" => false,
            # "params" => params,
            "error" => {
            }
        }

        if params[:user_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = User.where(id: params[:user_id]).first
            if user.nil?
                result["success"] = false
                result["error"]["message"] = "User doesn't exist!"
            else
                follow_requests = UserAction.where(to_user_id: params[:user_id], action_type: "follow_pending")
                if follow_requests.nil? || follow_requests.empty?
                    result["success"] = true
                    result["follow_requests"] = []
                    result["error"]["message"] = "No follow requests!"
                else
                    result["success"] = true
                    result["follow_requests"] = []
                    for request in follow_requests do 
                        request = request.attributes
                        request["user"] = User.where(id: request["from_user_id"]).first
                        result["follow_requests"].push(request)
                    end
                end
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
            to_user = User.where(id: params[:to_user_id]).first
            if user.nil? || to_user.nil?
                result["success"] = false
                result["error"]["message"] = "User[s] don't exist!"
            else
                block = UserAction.where(from_user_id: params[:to_user_id], to_user_id: params[:from_user_id], action_type: "block")
                if !block.nil? && !block.empty?
                    result["success"] = false
                    result["error"]["message"] = "This person has blocked you!"
                    result["block"] = block
                else
                    result["user"] = user
                    if params[:follow] === "true"
                        if to_user.protected
                            follow = UserAction.find_or_create_by(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow_pending")
                        else
                            follow = UserAction.find_or_create_by(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow")
                        end
                        # follow = user.build_user_action(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow")
                        follow.save
                        result["follow"] = follow
                        result["success"] = true
                    else
                        follow = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow").first
                        UserAction.delete(follow.id)
                        result["follow"] = false
                    end
                end
            end
        end
        render json: {result: result}
    end

    def approveFollowRequest
        result = {
            # params: params,
            "success" => false,
            "params" => params,
            "error" => {
            }
        }

        if params[:from_user_id].nil? || params[:to_user_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = User.where(id: params[:from_user_id]).first
            to_user = User.where(id: params[:to_user_id]).first
            if user.nil? || to_user.nil?
                result["success"] = false
                result["error"]["message"] = "User[s] don't exist!"
            else
                result["user"] = user
                follow = UserAction.find_by(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow_pending")
                if follow.nil?
                    result["success"] = false
                    result["error"]["message"] = "Follow request doesn't exist."
                else
                    result["success"] = true
                    follow.action_type = "follow"
                    follow.save
                    result["follow"] = follow
                end
            end
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
                follow_list = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow").
                        or(UserAction.where(from_user_id: params[:to_user_id], to_user_id: params[:from_user_id], action_type: "follow")).
                        or(UserAction.where(from_user_id: params[:to_user_id], to_user_id: params[:from_user_id], action_type: "follow_pending")).
                        or(UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "follow_pending")).
                        or(UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "mute")).
                        or(UserAction.where(from_user_id: params[:to_user_id], to_user_id: params[:from_user_id], action_type: "mute"))
                if !follow_list.nil? && !follow_list.empty? 
                    for follow in follow_list do
                        UserAction.delete(follow.id)
                    end
                end
                result["block"] = block
            else
                block = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "block").first
                if !block.nil?
                    UserAction.delete(block.id)
                end
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
            to_user = User.where(id: params[:to_user_id]).first
            if user.nil? || to_user.nil?
                result["success"] = false
                result["error"]["message"] = "User[s] don't exist!"
            else
                block = UserAction.where(from_user_id: params[:to_user_id], to_user_id: params[:from_user_id], action_type: "block")
                if !block.nil? && !block.empty?
                    result["success"] = false
                    result["error"]["message"] = "This person has blocked you!"
                    result["block"] = block
                else
                    user = User.where(id: params[:from_user_id]).first
                    result["user"] = user
                    if params[:mute] === "true"
                        mute = UserAction.find_or_create_by(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "mute")
                        mute.save
                        result["mute"] = mute
                    else
                        mute = UserAction.where(from_user_id: params[:from_user_id], to_user_id: params[:to_user_id], action_type: "mute").first
                        if !mute.nil?
                            UserAction.delete(mute.id)
                        end
                        result["mute"] = false
                    end
                    result["success"] = true
                end
            end
        end
        render json: {result: result}
    end

end
