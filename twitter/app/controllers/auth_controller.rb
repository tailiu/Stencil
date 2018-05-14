class AuthController < ApplicationController

    def new
        @new_user = User.new(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.build_credential(email: params[:email], password: params[:password])

        @result = {
            # params: params,
            "success" => false,
            "error" => {
            }
        }

        if @new_user.valid? && @new_credential.valid?
            @new_user.save
            @new_credential.save
            @result["success"] = true
            @result["user"] = @new_user
            reset_session
            session[:is_active] = true
            session[:user_id] = @new_user.id
        else 
            puts @new_user.errors.messages
            puts @new_credential.errors.messages

            if @new_user.errors.messages.key?(:handle) && @new_credential.errors.messages.key?(:email)
                @result["error"]["message"] = "Handle and Email already exist!"
            elsif @new_user.errors.messages.key?(:handle)
                @result["error"]["message"] = "Handle already exists!"
            elsif @new_credential.errors.messages.key?(:email)
                @result["error"]["message"] = "Email already exists!"
            else
                @result["error"]["message"] = "Invalid Credentials!"
            end

            @result["success"] = false
        end

        render json: {result: @result}
    end

    def verify

        @credentials = Credential.find_by(email: params[:email], password: params[:password])

        @result = {
            # params: params,
            "success" => false,
            "error" => {
            }
        }

        if @credentials != nil
            @result["success"] = true
            @result["user"]  = @credentials.user
            reset_session
            session[:is_active] = true
            session[:user_id] = @credentials.user.id
            session[:req_token] = rand(32**32).to_s(16)
            @result["session_id"] = session.id
            @result["req_token"] = session[:req_token]
            # @result["session"] = session
        else
            @result["success"] = false
            @result["error"]["message"] = "Invalid credentials!"
        end

        render json: {result: @result}
    end

    def logout
        session.clear
        reset_session
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        # if params[:user_id].nil?
        #     result["success"] = false
        #     result["error"]["message"] = "Incomplete params!"
        # else
        #     user = User.find_by_id(params[:user_id])
        #     if !user.nil?
        #         user.force_logout(session)
        #         reset_session
        #         result["success"] = true
        #     else
        #         result["success"] = false
        #         result["error"]["message"] = "User doesn't exist!"
        #     end
        # end
        render json: {result: result}
    end

    def is_logged_in
        result = {
            "params": params,
            "success" => false,
            "error" => {
            },
            "session" => session,
            "session_id" => session.id,
        }

        if params[:session_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            if session.id.to_s === params[:session_id].to_s
                if session[:is_active].nil?
                    result["success"] = false
                    result["error"]["message"] = "Session not found."
                else
                    if session[:is_active] == true
                        result["success"] = true
                    else
                        result["success"] = false
                        result["error"]["message"] = "Session not active."
                    end
                end
            else
                result["success"] = false
                result["error"]["message"] = "Session incompatible."
            end
        end
        render json: {result: result}
    end
    
end
