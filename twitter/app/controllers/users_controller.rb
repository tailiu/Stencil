class UsersController < ApplicationController
    def new
        @new_user = User.new(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.build_credential(email: params[:email], password: params[:password])

        @result = {
            # params: params,
            "success" => false,
            "error" => {
                "message": "",
            }
        }

        if @new_user.valid? && @new_credential.valid?
            @new_user.save
            @new_credential.save
            @result["success"] = true
            @result["user"] = @new_user
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
                "message": "",
            }
        }

        if @credentials.valid?
            @user = User.find(id: @credentials.user_id)
            if @user.valid?
                @result["success"] = true
                @result["user"]  = @user
            else
                puts @user.errors.messages
                @result["success"] = false
                @result["message"] = "User doesn't exist!"
            end
        else
            puts @credentials.errors.messages
            @result["success"] = false
            @result["message"] = "Invalid credentials!"
        end
        render json: @result
    end

end
